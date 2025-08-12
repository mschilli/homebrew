class Photogrep < Formula
    desc "Filter photos like Grep, using a Fyne GUI"
    homepage "https://github.com/mschilli/photogrep"
    url "https://github.com/mschilli/photogrep/archive/refs/tags/v0.02.tar.gz"
    license "Apache 2.0"

    depends_on "go" => :build

    test do
	assert_match "Usage", shell_output("#{bin}/ynabler --help", 1)
    end
end
