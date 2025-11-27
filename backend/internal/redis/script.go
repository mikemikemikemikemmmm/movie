package redis

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/lockSeat.lua
var lockSeatScriptSrc string
var lockSeatScript *redis.Script

//go:embed lua/releaseLock.lua
var releaseLockScriptSrc string
var releaseLockScript *redis.Script

//go:embed lua/reserve.lua
var reserveScriptSrc string
var reserveScript *redis.Script

//go:embed lua/check.lua
var checkScriptSrc string
var checkScript *redis.Script

func loadScripts() {
	lockSeatScript = redis.NewScript(lockSeatScriptSrc)
	releaseLockScript = redis.NewScript(releaseLockScriptSrc)
	reserveScript = redis.NewScript(reserveScriptSrc)
	checkScript = redis.NewScript(checkScriptSrc)
}
