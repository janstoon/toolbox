package kareless_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/janstoon/toolbox/bricks"
	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/kareless"
)

func TestLifeCycle(t *testing.T) {
	var (
		am avmod
		a1 app1
		a2 app2
		gw *commander1
	)

	k := kareless.Compile().
		Feed(settings{
			"port":   "123",
			"avstep": "5",
		}).
		Equip(func(_ *kareless.Settings, _ *kareless.InstrumentBank) []kareless.InstrumentCatalogue {
			return []kareless.InstrumentCatalogue{
				{
					Names: strings.Split("encryptor|decrypter", "|"),
					Builder: func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Instrument {
						am = newAvModifier(ss, ib)

						return am
					},
				},
			}
		}).
		Install(
			func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Application {
				a1 = newApp1(ss, ib)

				return a1
			},
			func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Application {
				a2 = newApp2(ss, ib)

				return a2
			},
		).
		Connect(func(ss *kareless.Settings, ib *kareless.InstrumentBank, apps []kareless.Application) kareless.Driver {
			adapter := struct {
				app1Commander1Adapter
				app2Commander1Adapter
			}{}
			for _, v := range apps {
				switch app := v.(type) {
				case app1:
					adapter.app1Commander1Adapter = newApp1Commander1Adapter(app)
				case app2:
					adapter.app2Commander1Adapter = newApp2Commander1Adapter(app)
				}
			}

			gw = newCommander1(ss, ib, adapter)

			return gw
		})

	started, stopped := make(chan bool), make(chan bool)
	k = k.
		AfterStart(
			func(ctx context.Context, ss *kareless.Settings, ib *kareless.InstrumentBank, apps []kareless.Application) error {
				started <- true

				go func() {
					<-ctx.Done()
					stopped <- true
				}()

				return nil
			},
		)

	assert.Nil(t, gw)
	assert.Nil(t, a1.encryptor)
	assert.Nil(t, a2.decrypter)
	assert.EqualValues(t, 0, am)

	ctx, stop := context.WithCancel(context.Background())
	go func() { _ = k.Run(ctx) }()

	select {
	case <-started:
	case <-time.After(500 * time.Millisecond):
		assert.Fail(t, "expected post hook to run")
	}
	assert.NotNil(t, a1.encryptor)
	assert.NotNil(t, a2.decrypter)
	assert.EqualValues(t, 5, am)
	assert.Equal(t, am, a2.decrypter)
	assert.Equal(t, am, a1.encryptor)

	assert.EqualValues(t, 123, gw.port)
	assert.NotNil(t, gw.app)
	assert.True(t, gw.isRunning())
	assert.Equal(t, "app2.bar app1.foo hello", gw.Bar(ctx, gw.Foo(ctx, "hello")))

	stop()
	select {
	case <-stopped:
	case <-time.After(500 * time.Millisecond):
		assert.Fail(t, "expected post hook context to get done")
	}
	assert.Eventually(t, func() bool { return !gw.isRunning() }, 500*time.Millisecond, 10*time.Millisecond)
}

type settings map[string]string

func (ss settings) Get(_ context.Context, key string) (any, error) {
	v, ok := ss[key]
	if ok {
		return v, nil
	}

	return nil, bricks.ErrNotFound
}

type message struct {
	payload string
}

type encryptor interface {
	Encrypt(src string) string
}

type app1 struct {
	encryptor encryptor
}

func newApp1(ss *kareless.Settings, ib *kareless.InstrumentBank) app1 {
	return app1{
		encryptor: kareless.ResolveInstrumentByType[encryptor](ib, "encryptor"),
	}
}

func (a app1) Foo(ctx context.Context, msg message) string {
	return a.encryptor.Encrypt(fmt.Sprintf("app1.foo %s", msg.payload))
}

type decrypter interface {
	Decrypt(src string) string
}

type app2 struct {
	decrypter decrypter
}

func newApp2(ss *kareless.Settings, ib *kareless.InstrumentBank) app2 {
	return app2{
		decrypter: kareless.ResolveInstrumentByType[decrypter](ib, "decrypter"),
	}
}

func (a app2) Bar(ctx context.Context, msg message) string {
	return fmt.Sprintf("app2.bar %s", a.decrypter.Decrypt(msg.payload))
}

type commander1App interface {
	Foo(ctx context.Context, msg string) string
	Bar(ctx context.Context, msg string) string
}

type commander1 struct {
	sync.RWMutex
	running bool
	port    int64

	app commander1App
}

func newCommander1(ss *kareless.Settings, ib *kareless.InstrumentBank, app commander1App) *commander1 {
	return &commander1{
		port: ss.GetInt64("port"),
		app:  app,
	}
}

func (c *commander1) Foo(ctx context.Context, msg string) string {
	return c.app.Foo(ctx, msg)
}

func (c *commander1) Bar(ctx context.Context, msg string) string {
	return c.app.Bar(ctx, msg)
}

func (c *commander1) Run(ctx context.Context) error {
	c.Lock()
	c.running = true
	c.Unlock()

	<-ctx.Done()
	c.Lock()
	c.running = false
	c.Unlock()

	return nil
}

func (c *commander1) isRunning() bool {
	c.RLock()
	defer c.RUnlock()

	return c.running
}

type app1Commander1Adapter struct {
	app app1
}

func newApp1Commander1Adapter(app app1) app1Commander1Adapter {
	return app1Commander1Adapter{
		app: app,
	}
}

func (adp app1Commander1Adapter) Foo(ctx context.Context, msg string) string {
	return adp.app.Foo(ctx, message{payload: msg})
}

type app2Commander1Adapter struct {
	app app2
}

func newApp2Commander1Adapter(app app2) app2Commander1Adapter {
	return app2Commander1Adapter{
		app: app,
	}
}

func (adp app2Commander1Adapter) Bar(ctx context.Context, msg string) string {
	return adp.app.Bar(ctx, message{payload: msg})
}

type avmod byte

func newAvModifier(ss *kareless.Settings, ib *kareless.InstrumentBank) avmod {
	return avmod(ss.GetByte("avstep"))
}

func (h avmod) Encrypt(src string) string {
	var dst strings.Builder
	for _, v := range src {
		dst.WriteByte(byte(v) + byte(h))
	}

	return dst.String()
}

func (h avmod) Decrypt(src string) string {
	var dst strings.Builder
	for _, v := range src {
		dst.WriteByte(byte(v) - byte(h))
	}

	return dst.String()
}
