cd dist
for i in *.apk; do
    [ -f "$i" ] || break
    cloudsmith push alpine --republish gsoc2/gsoc2-cli/alpine/any-version $i
done

for i in *.deb; do
    [ -f "$i" ] || break
    cloudsmith push deb --republish gsoc2/gsoc2-cli/any-distro/any-version $i
done

for i in *.rpm; do
    [ -f "$i" ] || break
    cloudsmith push rpm --republish gsoc2/gsoc2-cli/any-distro/any-version $i
done