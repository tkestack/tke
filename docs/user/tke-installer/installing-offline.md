1. Download the installation package with a network-connectable machine.

        version=v1.0.0 && wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-x86_64-$version.run{,.sha256}

2. Upload to a machine(A) in the target environment through tools such as rz or scp.
3. Run the installation package in the machine(A).

        version=v1.0.0 && sha256sum --check --status tke-installer-x86_64-$version.run.sha256 && chmod +x tke-installer-x86_64-$version.run && ./tke-installer-x86_64-$version.run

4. The remaining reference installation steps.