module github.com/emily33901/forgery

go 1.14

require (
	github.com/emily33901/forgery/core/events v0.0.0
	github.com/emily33901/forgery/core/filesystem v0.0.0
	github.com/emily33901/forgery/core/manager v0.0.0
	github.com/emily33901/forgery/core/materials v0.0.0-20200417143503-99942b2c0535
	github.com/emily33901/forgery/core/textures v0.0.0
	github.com/emily33901/forgery/core/world v0.0.0
	github.com/g3n/engine v0.1.1-0.20200214161420-db7282a2ba23
	github.com/galaco/KeyValues v1.4.1
	github.com/galaco/bsp v0.2.2 // indirect
	github.com/galaco/vmf v1.0.0 // indirect
	github.com/galaco/vpk2 v0.0.0-20181012095330-21e4d1f6c888 // indirect
	github.com/go-gl/gl v0.0.0-20190320180904-bf2b1f2f34d7
	github.com/go-gl/glfw v0.0.0-20200222043503-6f7a984d4dc4
	github.com/inkyblackness/imgui-go v1.12.0
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/emily33901/forgery/core/filesystem => ./core/filesystem/

replace github.com/emily33901/forgery/core/manager => ./core/manager/

replace github.com/emily33901/forgery/core/textures => ./core/textures/

replace github.com/emily33901/forgery/core/materials => ./core/materials/

replace github.com/emily33901/forgery/core/world => ./core/world/

replace github.com/emily33901/forgery/core/events => ./core/events/

replace github.com/inkyblackness/imgui-go => E:\src\gohack\github.com\inkyblackness\imgui-go
