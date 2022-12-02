{
  description = "app";
  outputs = { self, nixpkgs }:
    let
        system = "aarch64-linux";
        pkgs = nixpkgs.legacyPackages.${system};
    in
  {
      packages.${system} = {
          usb_proxy = import ./app.nix { inherit pkgs; };
      };
  };
}
