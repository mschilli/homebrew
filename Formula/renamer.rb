class Renamer < Formula
    git_url = "https://github.com/mschilli/homebrew"
    url git_url, :using => :git
    version "0.01"
    desc "Rename files in bulk"
    head git_url, :using => :git
    homepage "https://perlmeister.com"

    def install
      ENV['GOPATH'] = buildpath
      cd "projects/renamer" do
	  system "/usr/local/bin/go", "build"
	  bin.install "renamer"
      end
    end
end
