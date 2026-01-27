import { g as getPageContentFromS3 } from "../../../../chunks/s3.js";
const load = async ({ setHeaders }) => {
  setHeaders({
    "cache-control": "public, max-age=300, s-maxage=600, stale-while-revalidate=86400"
  });
  const content = await getPageContentFromS3("services");
  return {
    content,
    meta: {
      title: content.title || "Services",
      description: content.description || ""
    }
  };
};
export {
  load
};
