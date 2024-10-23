{
    
    description = "Kubernetes object analysis with recommendations for improved reliability and security";

    inputs.flake-utils.url = "github:numtide/flake-utils";

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem (system:
            let pkgs = import nixpkgs { inherit system; }; in
            rec {
                packages.kube-score = pkgs.buildGoModule {
                    name = "kube-score";
                    src = self;
                    vendorSha256 = "sha256-E9pcJsnoF/SKRCjrHZY8Ybd8kV1F3FYwdnLJ0mHyRLA="; # master
                    # vendorSha256 = "sha256-COY4AonAvJhH+io6Z7I9CsK1pnsK/Yi248QMkVPK6u0="; # v1.7.2
                    buildFlagsArray = ''
                      -ldflags=
                        -w -s
                        -X main.version=rolling
                        -X main.commit=${if self ? rev then self.rev else "dirty"}
                        -X main.date=${self.lastModifiedDate}
                        -X main.builtBy=nix
                    '';
                };
                defaultPackage = packages.kube-score;
                apps.kube-score = flake-utils.lib.mkApp { drv = packages.kube-score; };
                defaultApp = apps.kube-score;
            }
        );
}
