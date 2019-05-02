-- index to lookup member by provider ID
CREATE UNIQUE INDEX member_provider_idx ON member (provider, providerid);
