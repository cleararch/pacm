pkgname=pacm
pkgrel=1
pkgver=0.01
arch=(x86_64)
LICENSE="AGPLv3"
build(){
go build $srcdir/../main.go
}
package(){
	mkdir $pkgdir/usr
	mkdir $pkgdir/usr/bin
	mv main $pkgdir/usr/bin/pacm
}