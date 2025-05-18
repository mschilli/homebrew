class Ynabler < Formula
    desc "YNAB CSV conversion utilities"
    homepage "https://github.com/mschilli/go-ynabler"
    url "https://github.com/mschilli/go-ynabler/archive/refs/tags/v0.0.5.tar.gz"
    license "Apache 2.0"

    depends_on "go" => :build

    def install
	system "go", "build", "./cmd/ynabler"
	system "go", "build", "-o", bin/"ynabler-annotate", "annotate/cmd"
    end

    test do
	assert_match "Usage", shell_output("#{bin}/ynabler --help", 1)
	assert_match "Usage", shell_output("#{bin}/ynabler-annotate --help", 1)
    end
end
