import { useTranslation } from "react-i18next";
import Head from "next/head";

import { PersonalSettingsPage } from "@app/views/Settings/PersonalSettingsPage";

export default function PersonalSettings() {
  const { t } = useTranslation();

  return (
    <div className="bg-bunker-800 text-white h-full">
      <Head>
        <title>{t("common.head-title", { title: t("settings.personal.title") })}</title>
        <link rel="icon" href="/gsoc2.ico" />
      </Head>
      <PersonalSettingsPage />
    </div>
  );
}

PersonalSettings.requireAuth = true;
