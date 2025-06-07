# Open Issues

This document tracks features described in the README that are not yet implemented in the current code base.
Each section outlines an implementation plan.

## 1. LINE Messaging API Integration
The application currently logs weather results instead of sending LINE notifications.
### Implementation Plan
1. Add `github.com/line/line-bot-sdk-go/v7` to `go.mod` and vendor dependencies.
2. Create a service layer that wraps the LINE SDK for sending push messages.
3. Store `LINE_CHANNEL_ACCESS_TOKEN` and `LINE_CHANNEL_SECRET` in environment variables and pass them to the new service.
4. Modify `weatherUsecase.ProcessWeatherForUser` to send a message when `notify` is true.
5. Add unit tests using mocks for the LINE client.

## 2. LINE Login Authentication
The code only stores a LINE user ID and lacks authentication.
### Implementation Plan
1. Implement the LINE Login OAuth flow with callback endpoints.
2. Save access tokens and related info in the database (add migrations).
3. Refresh tokens when necessary and securely store them (e.g., encrypted at rest).
4. Add middleware to verify authenticated requests.
5. Provide tests for the new endpoints and token storage logic.

## 3. Automated Scheduling
Weather processing is triggered manually via `/api/process_weather`.
### Implementation Plan
1. Use Kubernetes CronJobs to fetch JMA weather data hourly and cache it.
2. Introduce a scheduler (Cron library or another CronJob) that runs `ProcessWeatherForUsersInTimeRange` for each user's preferred time.
3. Add the CronJob definitions to the Helm chart and ensure logs are captured.
4. Document how to run the scheduler locally using `make` or Docker Compose.

## 4. Redis Caching
Weather data is always fetched from JMA APIs in real time.
### Implementation Plan
1. Add a Redis instance to `compose.yml` and Helm `values.yaml`.
2. Use `github.com/go-redis/redis/v8` for caching weather responses.
3. Create a caching layer in the weather use case that checks Redis before making HTTP requests.
4. Write unit tests using `miniredis` to simulate Redis operations.

## 5. Chatbot Workflow for Settings
The README describes a chat-based UI for city selection and notification time.
### Implementation Plan
1. Implement webhook handlers using the LINE Messaging API to manage user interactions.
2. Present menus for selecting cities and notify times via flex messages or rich menus.
3. Persist user selections with the existing user use case.
4. Add tests for the conversation flow using mocked LINE events.

## 6. ArgoCD and CI/CD Enhancements
ArgoCD configuration is missing and the CI pipeline only runs Go tests.
### Implementation Plan
1. Add ArgoCD manifests under `build/kubernetes` for application and database deployments.
2. Extend GitHub Actions to build/push Docker images and trigger ArgoCD sync.
3. Document the deployment process in the README.

## 7. General Clean‑ups
Several configuration files reference outdated paths or omit scheduler resources.
### Implementation Plan
1. Update `compose.yml` to remove the non-existent `server` directory reference and build from the repository root.
2. Add CronJob templates to the Helm chart for hourly weather fetch and user notifications.
3. Review environment variable names and prune any that are unused.
