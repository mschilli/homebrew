class Ethtop < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.01"
    desc "Continuously display network interfaces in top format"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/ethtop" do
	  system "/usr/local/bin/go", "build", "ethtop.go", "ifconfig.go"
	  bin.install "ethtop"
      end
    end
end
