#!/bin/bash
get_major()
{
    echo $(cat VERSION | cut -d '.' -f1)
}

get_minor()
{
    echo $(cat VERSION | cut -d '.' -f2)
}

get_patch()
{
    echo $(cat VERSION | cut -d '.' -f3)
}

git_release()
{
    version=$(cat VERSION)
    git tag "v$version"
    git push origin "v$version"
}

git_push_version_bump()
{
    version=$(cat VERSION)
    git add VERSION
    git cm -m "Bumping VERSION file to $version"
    git push
}

if [ -z $1 ]; then
    echo "Missing version bump type. Options are major, minor, patch"
fi

if [ $1 == "major" ]; then
    git_release
    major=$(get_major)
    new_major=$(( major + 1 ))
    echo "$new_major.0.0" > VERSION
    git_push_version_bump
elif [ $1 == "minor" ]; then
    git_release
    major=$(get_major)
    minor=$(get_minor)
    new_minor=$(( minor + 1 ))
    echo "$major.$new_minor.0" > VERSION
    git_push_version_bump
elif [ $1 == "patch" ]; then
    git_release
    major=$(get_major)
    minor=$(get_minor)
    patch=$(get_patch)
    new_patch=$(( patch + 1 ))
    echo "$major.$minor.$new_patch" > VERSION
    git_push_version_bump
else
    echo "Invalid version bump type"
fi
