class Wifi < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.02"
    desc "Troubleshoot wifi network"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/wifi" do
	  system "go", "build"
	  bin.install "wifi"
      end
    end
end
