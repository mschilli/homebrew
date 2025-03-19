class Gittree < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.11"
    desc "Track git tree in Terminal UI"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/gittree" do
	  system "/usr/local/bin/go", "build"
	  bin.install "gittree"
      end
    end
end
