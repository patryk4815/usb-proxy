{
  pkgs ? import <nixpkgs> { }
}:
pkgs.buildGoModule {
  pname = "usb_proxy";
  version = "latest";

  src = ./.;

  nativeBuildInputs = [ pkgs.pkgconfig ];
  buildInputs = [ pkgs.libusb ];

#  vendorSha256 = "sha256-fSxtM+c6rkpKXALdaeMx3+gvAjtG2+mq0c+c9SvYzXY=";
#  vendorSha256 = pkgs.lib.fakeHash;
  vendorSha256 = "sha256-H5gbVcQdF5dgtaoLZ3DQJnUAmjjEXAqDLDm6GeNr/mw=";

  subPackages = [ "cmd/main" ];

  postInstall = ''
    mv $out/bin/main $out/bin/usb_proxy
  '';
#  meta = with lib; {
#    description = ''
#      Find security vulnerabilities, compliance issues, and infrastructure misconfigurations early in the development
#      cycle of your infrastructure-as-code with KICS by Checkmarx.
#    '';
#    homepage = "https://github.com/Checkmarx/kics";
#    license = licenses.asl20;
#    maintainers = with maintainers; [ patryk4815 ];
#  };
}