echo "Building package pidash_$VERSION-$ARCH.deb"

mkdir -p dist/deb || true

fpm -s dir -t deb -p dist/deb/pidash_$VERSION-$ARCH.deb \
        -n pidash -v $VERSION \
        --config-files /etc/pidash/dashboard.conf \
        --deb-priority optional --force \
        --deb-compression bzip2 \
        --description "PiDash" \
        -m "cats <admin@meow.tf>" --vendor "Meow.tf" \
        --before-install packaging/scripts/preinst.deb \
        --after-install packaging/scripts/postinst.deb \
        --url "https://meow.tf" \
        -a $ARCH \
        dist/$ARCH/pidash=/usr/bin/pidash \
        html=/var/lib/pidash \
        modules=/var/lib/pidash \
        modules/weather/module.example.conf=/etc/pidash/conf.d/weather.conf \
        packaging/pidash.service=/lib/systemd/system/pidash.service \
        dashboard.conf=/etc/pidash/dashboard.conf