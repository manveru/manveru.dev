let
  srcWithout = rootPath: ignoredPaths:
    let ignoreStrings = map (path: toString path) ignoredPaths;
    in builtins.filterSource
    (path: type: (builtins.all (i: i != path) ignoreStrings)) rootPath;

in import nixpkgsSource {
  config = { allowUnfree = true; };
  overlays = [
    (self: super: {
      infuse = super.stdenv.mkDerivation {
        pname = "infuse";
        version = "1.0.0.0";
        src = infuseSource;
        buildCommand = ''
          mkdir -p $out/bin
          cp $src $out/bin/infuse
          chmod +x $out/bin/infuse
        '';
      };

      go = self.go_1_12;

      inherit srcWithout;
      inherit (yarn2nix) yarn2nix mkYarnModules mkYarnPackage;
    })
  ];
}
