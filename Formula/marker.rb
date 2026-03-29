class Marker < Formula
    git_url = "https://github.com/mschilli/homebrew.git"
    url git_url, :using => :git
    version "0.04"
    desc "Text Highlighter"
    head git_url, :using => :git
    homepage "https://github.com/mschilli"

    depends_on "go" => :build

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/marker" do
          system "go", "build"
          bin.install "marker"
      end
    end
end
