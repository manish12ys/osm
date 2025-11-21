# Maintainer: Your Name <manishchowdary2006@gmail.com>
pkgname=omarchy-monitor
pkgver=1.0.0
pkgrel=1
pkgdesc="A beautiful, lightweight terminal-based system monitor built with Go"
arch=('x86_64' 'aarch64')
url="https://github.com/manish12ys/osm"
license=('MIT')
depends=()
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz::${url}/archive/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
    cd "${srcdir}/osm-${pkgver}"
    
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export CGO_LDFLAGS="${LDFLAGS}"
    export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
    
    go build -o ${pkgname} .
}

package() {
    cd "${srcdir}/osm-${pkgver}"
    
    # Install binary
    install -Dm755 "${pkgname}" "${pkgdir}/usr/bin/osm"
    
    # Install documentation
    install -Dm644 README.md "${pkgdir}/usr/share/doc/${pkgname}/README.md"
    
    # Install license if it exists
    if [ -f LICENSE ]; then
        install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
    fi
}
