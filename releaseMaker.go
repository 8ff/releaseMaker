package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func CreateNewRelease(client *github.Client, ctx context.Context, owner, repo string, newRelease *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	release, _, err := client.Repositories.CreateRelease(ctx, owner, repo, newRelease)
	if err != nil {
		// If an error occurs, delete the partially created release
		if release != nil {
			client.Repositories.DeleteRelease(ctx, owner, repo, *release.ID)
		}
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	return release, nil
}

func UploadAssetToReleaseByTag(client *github.Client, ctx context.Context, owner, repo, tagName, filePath, assetName string) error {
	// Get the existing release by tag name
	existingRelease, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tagName)
	if err != nil {
		return fmt.Errorf("failed to get release by tag: %w", err)
	}

	// Open the file you want to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Delete existing asset with the same name if it exists
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repo, *existingRelease.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to list release assets: %w", err)
	}

	for _, asset := range assets {
		if *asset.Name == assetName {
			_, err := client.Repositories.DeleteReleaseAsset(ctx, owner, repo, *asset.ID)
			if err != nil {
				return fmt.Errorf("failed to delete existing asset: %w", err)
			}
			break
		}
	}

	// Upload the file as an asset
	uploadOptions := &github.UploadOptions{
		Name: assetName, // Name of the file in the release
	}
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, owner, repo, *existingRelease.ID, uploadOptions, file)
	if err != nil {
		return fmt.Errorf("failed to upload asset: %w", err)
	}

	return nil
}

func DeleteReleaseByTag(client *github.Client, ctx context.Context, owner, repo, tagName string) error {
	// Get the existing release by tag name
	existingRelease, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tagName)
	if err != nil {
		return fmt.Errorf("failed to get release by tag: %w", err)
	}

	// Delete the release
	_, err = client.Repositories.DeleteRelease(ctx, owner, repo, *existingRelease.ID)
	if err != nil {
		return fmt.Errorf("failed to delete release: %w", err)
	}

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [command] [arguments]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  create [owner/repo] [tag] [name] [body] - Create a new release")
	fmt.Fprintln(os.Stderr, "  upload [owner/repo] [tag] [file] [assetName] - Upload a file as an asset to an existing release")
	fmt.Fprintln(os.Stderr, "  delete [owner/repo] [tag] - Delete an existing release")
}

func createClient(token string) *github.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30-second timeout
	defer cancel()                                                           // Make sure to cancel the context when done to release resources
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		printUsage()
		os.Exit(2)
	}

	// Read TOKEN from env
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN is not set")
		os.Exit(2)
	}

	switch args[0] {
	case "create":
		if len(args) < 6 {
			fmt.Fprintf(os.Stderr, "Usage: %s create [owner/repo] [tag] [name] [body]\n", os.Args[0])
			os.Exit(2)
		}
		ownerRepo := strings.Split(args[1], "/")
		if len(ownerRepo) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid owner/repo argument: %s\n", args[1])
			os.Exit(2)
		}
		owner := ownerRepo[0]
		repo := ownerRepo[1]
		// Create an authenticated client
		client := createClient(token)

		// Read release tag, name, body from args
		tagName := args[2]
		releaseName := args[3]
		releaseBody := args[4]

		// Check for empty args
		if owner == "" || repo == "" || tagName == "" || releaseName == "" || releaseBody == "" {
			fmt.Fprintf(os.Stderr, "Invalid arguments: %s\n", args[1:])
			os.Exit(2)
		}

		// Define the new release information
		newRelease := &github.RepositoryRelease{
			TagName:    github.String(tagName),
			Name:       github.String(releaseName),
			Body:       github.String(releaseBody),
			Prerelease: github.Bool(false),
		}

		// Create a new release
		_, err := CreateNewRelease(client, context.Background(), owner, repo, newRelease)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Release created successfully!")
		os.Exit(0)
	case "upload":
		if len(args) < 5 {
			fmt.Fprintf(os.Stderr, "Usage: %s upload [owner/repo] [tag] [file]\n", os.Args[0])
			os.Exit(2)
		}
		ownerRepo := strings.Split(args[1], "/")
		if len(ownerRepo) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid owner/repo argument: %s\n", args[1])
			os.Exit(2)
		}
		owner := ownerRepo[0]
		repo := ownerRepo[1]
		// Create an authenticated client
		client := createClient(token)

		// Read release tag, name, body from args
		tagName := args[2]
		filePath := args[3]
		assetName := args[4]

		// Check for empty args
		if owner == "" || repo == "" || tagName == "" || filePath == "" || assetName == "" {
			fmt.Fprintf(os.Stderr, "Invalid arguments: %s\n", args[1:])
			os.Exit(2)
		}

		// Upload the file as an asset
		err := UploadAssetToReleaseByTag(client, context.Background(), owner, repo, tagName, filePath, assetName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Release uploaded successfully!")
		os.Exit(0)
	case "delete":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s delete [owner/repo] [tag]\n", os.Args[0])
			os.Exit(2)
		}
		ownerRepo := strings.Split(args[1], "/")
		if len(ownerRepo) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid owner/repo argument: %s\n", args[1])
			os.Exit(2)
		}
		owner := ownerRepo[0]
		repo := ownerRepo[1]
		// Create an authenticated client
		client := createClient(token)

		// Read release tag, name, body from args
		tagName := args[2]

		// Check for empty args
		if owner == "" || repo == "" || tagName == "" {
			fmt.Fprintf(os.Stderr, "Invalid arguments: %s\n", args[1:])
			os.Exit(2)
		}

		// Delete the release
		err := DeleteReleaseByTag(client, context.Background(), owner, repo, tagName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Release deleted successfully!")
		os.Exit(0)
	default:
		printUsage()
		os.Exit(2)
	}
}
