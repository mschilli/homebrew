class Gitmeta < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.10"
    desc "Parse .gmf git meta format files and update local repos"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/gitmeta" do
	  system "/usr/local/bin/go", "build"
	  bin.install "gitmeta"
      end
    end
end
