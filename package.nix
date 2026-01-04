{
  buildGoModule,
  buildPackages,
  fetchFromGitHub,
  installShellFiles,
  lib,
  stdenv,
}:

buildGoModule (finalAttrs: {
  pname = "keys";
  version = "2026.1.15";

  src = fetchFromGitHub {
    owner = "adreasnow";
    repo = "keys";
    tag = "v${finalAttrs.version}";
    hash = "sha256-IJ0fgnjn9xIshfFnr7AX15V4MCFnn9XcqzKUbUJmSAc=";
  };
  vendorHash = "sha256-diCC3gA2hnAlzH3E7syMrKD3yebU+ZsSnicWg8ZW6x8=";

  ldflags = [ "-s -w -X main.version=${finalAttrs.version}" ];

  nativeBuildInputs = [ installShellFiles ];

  subPackages = [ "cmd/keys" ];

  postInstall =
    let
      keysBin =
        if stdenv.buildPlatform.canExecute stdenv.hostPlatform then
          "$out"
        else
          lib.getBin buildPackages.keys;
    in
    ''
      installShellCompletion --cmd keys \
        --bash <(${keysBin}/bin/keys completion bash) \
        --fish <(${keysBin}/bin/keys completion fish) \
        --zsh <(${keysBin}/bin/keys completion zsh)
    '';

  meta = {
    description = "Lightweight Go wrapper around your OS's keychain to act as a simple CLI tool for managing secrets.";
    homepage = "https://github.com/adreasnow/keys";
    mainProgram = "keys";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [
      adreasnow
    ];
  };
})
