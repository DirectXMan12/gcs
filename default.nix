with import <nixpkgs> {};

let
    src = fetchFromGitHub {
        owner = "directxman12";
        repo = "gcs";
        rev = "v5.8.0-custom.0.0";
        sha256 = "026vi3m52qfvaq1d8m2h4qwh76f7nr1s9y8gghkxi2sdgk42y69z";
    };
    unstable = import <unstable> {};
in

stdenv.mkDerivation {
    name = "gcs-5.8.0-custom.0.0";

    inherit src;

    buildInputs = [
        gcc
        unstable.go
        zlib
        xorg.libXext
        xorg.libX11
        xorg.libXrender
        xorg.libXtst
        xorg.libXi
        xorg.libXcursor
        xorg.libXrandr
        xorg.libXinerama
        xorg.libXxf86vm
        libGL
        wayland
        wayland-protocols
        libxkbcommon
        freetype
        pkg-config
        fontconfig
    ];

    buildPhase = ''
        ./build.sh
    '';

    installPhase = ''
        mkdir -p $out/bin
        cp gcs $out/bin/gcs
        chmod +x $out/bin/gcs
    '';

    # NB(directxman12): this doesn't cover the library, which is managed by GCS itself
}
