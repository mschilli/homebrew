class Ynabler < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.0.1"
    desc ""
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/ynabler" do
	  system "/usr/local/bin/go", "build"
	  bin.install "ynabler"
      end
    end
end
