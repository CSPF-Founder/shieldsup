import argparse
from bson.objectid import ObjectId

from app.modules.reporter import run_reporter


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-t", "--target-id", dest="target_id", required=True)
    parser.add_argument(
        "-m",
        "--module",
        help="Module to run",
        dest="module",
        required=True,
        choices=["reporter"],
    )
    args = parser.parse_args()
    if ObjectId.is_valid(args.target_id):
        # run_scan_for_id(args.target_id)
        if args.module == "reporter":
            target_id = ObjectId(args.target_id)
            run_reporter(target_id=target_id)
    else:
        print("Invalid object id received")
        exit(1)


if __name__ == "__main__":
    main()
