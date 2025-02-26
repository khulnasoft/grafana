import { useState, useEffect } from 'react';
import { Button, Box, Text, Stack } from '@grafana/ui';
import { Feature, feature_ini, ngrok_example, root_url_ini } from './types';
import { SetupModal } from './SetupModal';
import { FeatureCard } from './FeatureCard';
import { getConfigurationStatus } from './utils';

export const FeatureList = () => {
  const [features, setFeatures] = useState<Feature[]>([]);
  const [showInstructionsModal, setShowInstructionsModal] = useState(false);
  const [activeFeature, setActiveFeature] = useState<Feature | null>(null);

  const { hasPublicAccess, hasImageRenderer, hasRequiredFeatures } = getConfigurationStatus();

  useEffect(() => {
    // Initialize features with their current status
    const featuresList: Feature[] = [
      {
        title: 'Provision As-Code',
        description: 'Provision your dashboards from Github or other storage system',
        additional: false,
        setupSteps: [],
        isConfigured: hasRequiredFeatures,
        icon: 'sync',
      },
      {
        title: 'Pull Request Collaboration',
        description: 'Collaborate with your team using pull requests',
        additional: false,
        setupSteps: [],
        icon: 'code-branch',
        isConfigured: hasRequiredFeatures,
      },
      {
        title: 'Migrate Your Dashboards',
        description: 'Migrate your dashboards to Github or other storage system',
        additional: false,
        setupSteps: [],
        icon: 'cloud-upload',
        isConfigured: hasRequiredFeatures,
      },
      {
        title: 'Github Webhooks',
        description: 'Seamless Github provisioning and collaboration with pull requests',
        additional: true,
        icon: 'github',
        isConfigured: hasPublicAccess && hasRequiredFeatures,
        setupSteps: [
          {
            title: 'Start ngrok for temporary public access',
            description: 'Run this command to create a secure tunnel to your local Grafana:',
            code: 'ngrok http 3000',
          },
          {
            title: 'Copy your public URL',
            description: 'From the ngrok output, copy the https:// forwarding URL that looks like this:',
            code: ngrok_example,
            copyCode: false,
          },
          {
            title: 'Update your Grafana configuration',
            description: 'Add this to your custom.ini file, replacing the URL with your actual ngrok URL:',
            code: root_url_ini,
          },
        ],
      },
      {
        title: 'Preview Snapshots',
        description: 'Attach preview images to pull requests comments',
        icon: 'camera',
        additional: true,
        isConfigured: hasImageRenderer && hasPublicAccess && hasRequiredFeatures,
        setupSteps: [
          {
            title: 'Install Node.js',
            description: 'Install Node.js 16 or later on your system',
          },
          {
            title: 'Clone the Image Renderer Repository',
            description: 'Clone the renderer repository:',
            code: 'git clone https://github.com/grafana/grafana-image-renderer.git',
          },
          {
            title: 'Build the Renderer',
            description: 'Navigate to the directory and build:',
            code: 'cd grafana-image-renderer\nnpm install\nnpm run build',
          },
          {
            title: 'Run the Renderer Service',
            description: 'Start the renderer service:',
            code: 'node build/app.js server --port=8081',
          },
          {
            title: 'Configure Grafana',
            description: 'Add these settings to your grafana.ini file:',
            code: `[rendering]
rendering_server_url = http://localhost:8081/render
rendering_callback_url = http://your-grafana-instance/`,
          },
        ],
      },
    ];

    setFeatures(featuresList);
  }, []);

  // Add a state variable to store the basic setup
  const [basicSetup] = useState<Feature>({
    title: 'Provisioning',
    description: 'Enable required Grafana features for provisioning',
    additional: false,
    icon: 'cog',
    isConfigured: hasRequiredFeatures,
    setupSteps: [
      {
        title: 'Enable Required Feature Toggles',
        description: 'Add these settings to your custom.ini file to enable necessary features:',
        code: feature_ini,
      },
    ],
  });

  const handleShowInstructions = (feature: Feature) => {
    // only show modal if feature is not configured
    if (!feature.isConfigured) {
      setActiveFeature(feature);
      setShowInstructionsModal(true);
    }
  };

  const handleInstructionsClose = () => {
    setShowInstructionsModal(false);
    setActiveFeature(null);
  };

  // Separate required and optional features
  const requiredFeatures = features.filter((feature) => !feature.additional);
  const optionalFeatures = features.filter((feature) => feature.additional);

  return (
    <Stack direction="column" gap={4}>
      <Text element="h1" variant="h2" textAlignment="center">
        All Features
      </Text>
      <Stack direction="row" gap={2} justifyContent="center">
        {requiredFeatures.map((feature, index) => (
          <FeatureCard
            key={index}
            feature={feature}
            onSetup={() => handleShowInstructions(basicSetup)}
            showSetupButton={true}
          />
        ))}
        {optionalFeatures.map((feature, index) => (
          <FeatureCard
            key={index}
            feature={feature}
            onSetup={() => handleShowInstructions(feature)}
            showSetupButton={true}
          />
        ))}
      </Stack>

      {showInstructionsModal && activeFeature && (
        <SetupModal feature={activeFeature} isOpen={true} onDismiss={handleInstructionsClose} />
      )}
    </Stack>
  );
};
