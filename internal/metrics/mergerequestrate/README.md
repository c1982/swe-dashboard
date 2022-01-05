# What is MR Rate?

> https://about.gitlab.com/blog/2020/08/27/measuring-engineering-productivity-at-gitlab/#what-is-mr-rate

We use this metric because:

    We want to incentivize everyone to iterate and break down work into smaller MRs because smaller MRs have a faster review time and get merged faster (better developer and maintainer review experience)
    The quicker we can deliver features to users, the faster we can iterate upon them
    Every MR into the codebase improves the codebase, and every improvement has the downstream effect of making the product better

When viewed at an organization level, this metric helps us understand how productivity in the organization changes over time. Although this metric seems simple, it actually requires a lot of detailed analysis as there are many situations to examine:

    New team vs. established team
    Team performance issues (blocking work or incorrect iteration work breakdown)
    Individual growth (and performance management)
    Community contributions vs. independent team contributions
    Operational productivity constraints

At first, we measured MRs based on labels associated with the product domain (which generally maps to an existing engineering team). As an open core company, this allowed us to easily aggregate community contributions into the metric. We wanted to account for them because we want to continue encouraging team members to support community contribution MRs and recognize that these MRs continue to help provide the product with more value to users.

Unfortunately, as our organization grew over time, this metric became confusing. Although we had a bot that would label MRs, we occasionally had bad data and mislabeled MRs. In addition, certain teams with product areas that were more mature had more community contributions than others. The combination of these issues evolved the metric into multiple types.

    MR Rate measured through labeling
    Team MR Rate measured through MR authorship (also known as Narrow MR Rate)

It's likely that over time this may continue to evolve but for now, these new types of MR Rates have brought more clarity within our organization.

## links:

* https://about.gitlab.com/blog/2020/08/27/measuring-engineering-productivity-at-gitlab/#what-is-mr-rate
* https://www.linkedin.com/pulse/size-pullmerge-request-more-important-than-you-think-rodrigo-miguel/
* https://stiltsoft.com/blog/2021/03/pull-request-analytics-how-to-visualize-cycle-time-lead-time-and-get-insights-for-improvement/