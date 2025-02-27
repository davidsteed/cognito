module.exports = {
  backoffice: {
    output: {
      mode: "tags",
      target: "./endpoints/backOfficeFromFileSpec.ts",
      schemas: "./model",
      client: "react-query",
      workspace: "./app/generator",
      prettier: true,
      override: {
        mutator: {
          path: "./mutator/customInstance.ts",
          name: "customInstance",
        },
        query: {
          useQuery: true,
        },
      },
    },
    input: {
      target: "./api/location.yaml",
    },
    hooks: {
      afterAllFilesWrite: "prettier . --write",
    },
  },
};
