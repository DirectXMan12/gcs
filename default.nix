with import <nixpkgs> {};

let
    src = fetchFromGitHub {
		owner = "directxman12";
		repo = "gcs";
		rev = "v4.37.1-custom.0.0";
        sha256 = "026vi3m52qfvaq1d8m2h4qwh76f7nr1s9y8gghkxi2sdgk42y69z";
    };
in

stdenv.mkDerivation {
    name = "gcs-4.37.1-custom.0.0";

    inherit src;

    buildInputs = [
      jdk
      zlib
      xorg.libXext
      xorg.libX11
      xorg.libXrender
      xorg.libXtst
      xorg.libXi
      freetype
	];

	buildPhase = ''
		./bundle.sh -u
	'';

    wrapperLauncher = ''
        #!${stdenv.shell}
		export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:${zlib}/lib:${xorg.libXext}/lib:${xorg.libX11}/lib:${xorg.libXrender}/lib:${xorg.libXtst}/lib:${xorg.libXi}/lib:${freetype}/lib";
		/share/gcs/bin/GCS "$@"
    '';

	# gcs seems to want this weird layout
    installPhase = ''
		mkdir -p $out/bin
		mkdir -p $out/share/gcs/{bin,lib}

		cp GCS/bin/GCS $out/share/gcs/bin/GCS
		cp -r GCS/lib/{GCS.png,libapplauncher.so,app,runtime} $out/share/gcs/lib
        echo "$wrapperLauncher" > $out/bin/gcs
        substituteInPlace $out/bin/gcs --replace /share/gcs/bin/GCS $out/share/gcs/bin/GCS
        chmod +x $out/bin/gcs
    '';

	# NB(directxman12): this doesn't cover the library, which is managed by GCS itself
}
