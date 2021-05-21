#!/usr/bin/env python

# MIT License
#
# Copyright (c) 2021 Theodor Amariucai
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

import argparse
import cpustats
import transfer
import burstiness
import imgsize


def parse_args():
    parser = argparse.ArgumentParser(description='Visualize experiment results.')
    parser.add_argument('--type', type=str, default='cpustats',
                        help='Experiment type to be visualized.')  # [cpustats, cdfs, transfer, imgsize, imgsize-burstiness] are the possible options
    parser.add_argument('--path', type=str, default='.', help='Path where the experiment results are located.')
    return parser.parse_args()


# Examples:
# 1. CPU stats:
# python plot.py --type cpustats --path ../latency-samples/cloudlab/cpu-stats/AWS/
# 2. Transfers:
# python plot.py --type transfer --path ../latency-samples/cloudlab/data-transfer/AWS/inline/
# 3. Burstiness CDFs
# python plot.py --type burstiness --path ../latency-samples/cloudlab/burstiness/AWS/Tuesday\,\ 16-Mar-21\ 08\:19\:39\ MDT/
# 4. Cold starts image size:
# python plot.py --type imgsize --path ../latency-samples/cloudlab/image-size/AWS/1536MB\ memory\,\ st1s\ imgsize\ experiment/

if __name__ == "__main__":
    args = parse_args()
    print(f"Path is {args.path} and visualization type is {args.type}")

    if "AWS" in args.path:
        args.provider = "AWS"
    elif "vhive" in args.path or "vHive" in args.path:
        args.provider = "vHive"
    else:
        raise Exception(f"Unrecognized provider in path {args.path}")
    print(f'Identified provider is {args.provider}')

    if args.type == "cpustats":
        cpustats.plot_cpu_stats(args)
    elif args.type == "transfer":
        transfer.plot_data_transfer_stats(args)
    elif args.type == "burstiness":
        burstiness.plot_cdfs(args)
    elif args.type == "imgsize":
        imgsize.plot_imgsize_stats(args)
    else:
        raise Exception(f"Unsupported visualization type {args.type}")
