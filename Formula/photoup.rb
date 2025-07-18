class Photoup < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.01"
    desc "Prepare a stack of photos for public viewing"
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/photoup" do
          system "go", "build"
          bin.install "photoup"
      end
    end
end
