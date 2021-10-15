package main

type tomlConfig struct {
	Title      string
	Connection connection `toml:"connection"`
}

type connection struct {
	Port int
}
