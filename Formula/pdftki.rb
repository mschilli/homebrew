class Pdftki < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.1"
    desc "Convenience wrapper around pdftk, -e pops up an editor to edit the command line"
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/pdftki" do
          system "go", "build"
          bin.install "pdftki"
      end
    end
end
