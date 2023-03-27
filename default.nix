with import <nixpkgs> {};

let
    src = fetchFromGitHub {
        owner = "directxman12";
        repo = "gcs";
        rev = "v5.8.0-custom.0.0";
		hash = "sha256-BPXLLfLwSwxvLLFDdQJzTqeHIGvi5U+BEJuthzxDQe8=";
    };
    unstable = import <unstable> {};
    libs = [
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
        fontconfig
    ];
in

buildGo119Module {
    name = "gcs-5.8.0-custom.0.0";

    inherit src;

	# go mod vendor doesn't include C code, so for richardwilkes/pdf
	# we need to proxy instead
	proxyVendor = true;

    nativeBuildInputs = [
        bash
        gcc
        pkg-config
    ] ++ libs;

    buildInputs = [
    ] ++ libs;

    vendorSha256 = "sha256-JldsJTdXUUpY4c0HTOd35RcVznawvAtJ9Su70BqkKMk=";
    
    ldflags="-X github.com/richardwilkes/toolbox/cmdline.AppVersion=5.8.0";

    # NB(directxman12): this doesn't cover the library, which is managed by GCS itself
}
