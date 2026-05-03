class Configpatch < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.0.4"
    desc "Command line utility to patch config files, based on https://github.com/mschilli/go-configpatch"
    head git_url, :using => :git
    homepage "https://github.com/mschilli/go-configpatch"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/configpatch" do
          system "go", "build"
          bin.install "configpatch"
      end
    end
end
