import { Fragment } from "react";
import { Tab } from "@headlessui/react";

import { OrgPermissionActions, OrgPermissionSubjects } from "@app/context";
import { withPermission } from "@app/hoc";

import { BillingCloudTab } from "../BillingCloudTab";
import { BillingDetailsTab } from "../BillingDetailsTab";
import { BillingReceiptsTab } from "../BillingReceiptsTab";
import { BillingSelfHostedTab } from "../BillingSelfHostedTab";

const tabs = [
  { name: "Gsoc2 Cloud", key: "tab-gsoc2-cloud" },
  { name: "Gsoc2 Self-Hosted", key: "tab-gsoc2-self-hosted" },
  { name: "Receipts", key: "tab-receipts" },
  { name: "Billing details", key: "tab-billing-details" }
];

export const BillingTabGroup = withPermission(
  () => {
    return (
      <Tab.Group>
        <Tab.List className="mt-8 mb-6 border-b-2 border-mineshaft-800">
          {tabs.map((tab) => (
            <Tab as={Fragment} key={tab.key}>
              {({ selected }) => (
                <button
                  type="button"
                  className={`w-30 py-2 mx-2 mr-4 font-medium text-sm outline-none ${
                    selected ? "border-b border-white text-white" : "text-mineshaft-400"
                  }`}
                >
                  {tab.name}
                </button>
              )}
            </Tab>
          ))}
        </Tab.List>
        <Tab.Panels>
          <Tab.Panel>
            <BillingCloudTab />
          </Tab.Panel>
          <Tab.Panel>
            <BillingSelfHostedTab />
          </Tab.Panel>
          <Tab.Panel>
            <BillingReceiptsTab />
          </Tab.Panel>
          <Tab.Panel>
            <BillingDetailsTab />
          </Tab.Panel>
        </Tab.Panels>
      </Tab.Group>
    );
  },
  { action: OrgPermissionActions.Read, subject: OrgPermissionSubjects.Billing }
);
