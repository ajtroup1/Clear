package modules

import "github.com/ajtroup1/clear/object"

func Register(env *object.Environment) {
	env.SetModule("math", MathBuiltins)
	env.SetModule("strings", StringsBuiltins)
	env.SetModule("arrays", ArraysBuiltins)
	// env.SetModule("rand", RandBuiltins)
	env.SetModule("io", IOBuiltins)
	// env.SetModule("os", OSBuiltins)
	// env.SetModule("time", TimeBuiltins)

}
