1. Download the installation package with a network-connectable machine.

        # Choose installation package according to your installation node CPU architecture [amd64, arm64]
        arch=amd64 version=v1.3.0 && wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-linux-$arch-$version.run{,.sha256}

2. Upload to a machine(A) in the target environment through tools such as rz or scp.
3. Run the installation package in the machine(A).

        # Choose installation package according to your installation node CPU architecture [amd64, arm64]
        arch=amd64 version=v1.3.0 && sha256sum --check --status tke-installer-linux-$arch-$version.run.sha256 && chmod +x tke-installer-linux-$arch-$version.run && ./tke-installer-linux-$arch-$version.run

4. The remaining reference installation steps.