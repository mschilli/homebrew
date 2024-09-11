class Murmur < Formula
    git_url = "https://github.com/mschilli/murmur"
    url git_url, :using => :git
    version "0.0.1"
    desc "Handle secrets with Murmur"
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/murmur" do
	  system "go", "build"
	  bin.install "murmur-fill"
      end
    end
end
