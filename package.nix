{
  fetchFromGitHub,
  lib,
  go,
  stdenv,
}:
let
  pname = "auto-pull";
  appName = "Auto Pull";
  version = "2025.12.8";
  meta = {
    description = "Auto GitHub Puller";
    mainProgram = "autopull";
    maintainers = with lib.maintainers; [
      adreasnow
    ];
  };

  src = fetchFromGitHub {
    owner = "adreasnow";
    repo = "auto-pull";
    tag = "v${version}";
    hash = "sha256-x2LXvANLQssKoLXjhqoglFaJoWlZDJjCPpJb0pPMOmU=";
  };

  darwin = stdenv.mkDerivation {
    inherit
      pname
      version
      src
      ;
    meta = meta // {
      platforms = [
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    };

    nativeBuildInputs = [
      go
    ];

    buildPhase = ''
      runHook preBuild

      export HOME=$TMPDIR
      export GOCACHE=$TMPDIR/go-cache
      export GOMODCACHE=$TMPDIR/go-mod

      # Create app bundle structure
      mkdir -p "${appName}.app/Contents/MacOS"
      mkdir -p "${appName}.app/Contents/Resources"

      # Build the binary
      go build -o "${appName}.app/Contents/MacOS/autopull" ./cmd

      # Create Info.plist
      cat > "${appName}.app/Contents/Info.plist" << EOF
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
        <key>CFBundleExecutable</key>
        <string>autopull</string>
        <key>CFBundleIconFile</key>
        <string>icon.icns</string>
        <key>CFBundleGetInfoString</key>
        <string>${appName}</string>
        <key>CFBundleIdentifier</key>
        <string>com.github.adreasnow.auto-pull</string>
        <key>CFBundleName</key>
        <string>${appName}</string>
        <key>CFBundleShortVersionString</key>
        <string>${version}</string>
        <key>CFBundleInfoDictionaryVersion</key>
        <string>6.0</string>
        <key>CFBundlePackageType</key>
        <string>APPL</string>
        <key>IFMajorVersion</key>
        <integer>1</integer>
        <key>IFMinorVersion</key>
        <integer>0</integer>
        <key>NSHighResolutionCapable</key><true/>
        <key>NSSupportsAutomaticGraphicsSwitching</key><true/>
      </dict>
      </plist>
      EOF

      cp icons/icon.icns "${appName}.app/Contents/Resources/"
      cp icons/cloud.png "${appName}.app/Contents/Resources/"
      cp icons/sun.png "${appName}.app/Contents/Resources/"
      cp icons/warning.png "${appName}.app/Contents/Resources/"

      runHook postBuild
    '';

    installPhase = ''
      runHook preInstall
      mkdir -p "$out/Applications"
      cp -r "${appName}.app" "$out/Applications/"
      runHook postInstall
    '';
  };
in
darwin
