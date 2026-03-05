{
  lib,
  buildGoModule,
  go_1_25,
}:

let
  manifest = builtins.fromJSON (builtins.readFile ../.release-please-manifest.json);
in
buildGoModule rec {
  pname = "stl";
  version = manifest.".";

  vendorHash = "sha256-OWfYAhV9jFuRTMtjjMdut2htFcYzjDl/RqgzesOBptk=";

  src = ../.;

  nativeBuildInputs = [ go_1_25 ];

  doCheck = false;

  meta = {
    description = "Stainless CLI";
    homepage = "https://github.com/stainless-api/stainless-api-cli";
    license = lib.licenses.asl20;
    mainProgram = pname;
  };
}
