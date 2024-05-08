class Ynabler < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.0.3"
    desc ""
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/ynabler" do
	  system "go", "build"
	  bin.install "ynabler"
      end
    end
end
