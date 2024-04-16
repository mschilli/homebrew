class Autorot < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.02"
    desc "Rotate images according to their EXIF info"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/autorot" do
	  system "/usr/local/bin/go", "build"
	  bin.install "autorot"
      end
    end
end
