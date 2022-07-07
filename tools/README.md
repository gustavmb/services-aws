# AWS Service Engineer Tools

The intention of this space is to share micro tools that can help us in our daily job duties when working with Virtualization - focusing on KVM.

## Pinning logical CPUS through Unix Domain Sockets approach

This is a Go implementation that sets affinity to specific CPUS given a JSON file located at sockdata directory.

```
{
  "VM": [
    {
      "sockpath": "/path/to/unix_1.sock",
      "vcpus": [0,1,2,3]
    },
    {
      "sockpath": "/path/to/unix_2.sock",
      "vcpus": [4,5,6,7]
    },
    {
      "sockpath": "/path/to/unix_3.sock",
      "vcpus": [8,9,10,11]
    },
    {
      "sockpath": "/path/to/unix_4.sock",
      "vcpus": [12,13,14,15]
    },
    {
      "sockpath": "/path/to/unix_5.sock",
      "vcpus": [16,17,18,19]
    }
  ]
}

```
The figure shown above describes the location of the socket file after launching the VM - sockpath - and defines the vcpus that we want to use to PIN to the host.

### How to Run

A binary is provided in this repo "sock_kvm" and you can run it out of the box, you just need to run the binary after defining the affinity and sockpath in the file described above. You can always compile the code if you want to make some modifications according with your necessities.

```
./sock_kvm
```
