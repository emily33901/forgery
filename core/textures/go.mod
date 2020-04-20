module github.com/emily33901/forgery/core/textures

go 1.14

require (
	github.com/emily33901/forgery/core/filesystem v0.0.0-20200414124732-e68de349b69b
	github.com/emily33901/forgery/core/manager v0.0.0-20200417100142-4fbd8606393a
	github.com/emily33901/vtf v0.0.0-20190613094935-fc3c5b74c85e
)

// replace github.com/emily33901/forgery/core/filesystem => ../core/filesystem/
replace github.com/emily33901/forgery/core/manager => ../manager/

replace github.com/emily33901/forgery/core/filesystem => ../filesystem/
