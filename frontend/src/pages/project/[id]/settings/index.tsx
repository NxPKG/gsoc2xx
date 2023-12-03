import { useTranslation } from "react-i18next";
import Head from "next/head";

import { ProjectSettingsPage } from "@app/views/Settings/ProjectSettingsPage";

const ProjectSettings = () => {
  const { t } = useTranslation();

  return (
    <>
      <Head>
        <title>{t("common.head-title", { title: t("settings.project.title") })}</title>
        <link rel="icon" href="/gsoc2.ico" />
      </Head>
      <ProjectSettingsPage />
    </>
  );
};

export default ProjectSettings;

ProjectSettings.requireAuth = true;
