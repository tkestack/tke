import * as React from 'react';
export function NavigateLink({
  title,
  router,
  fragment,
  queries,
  children,
  onClick,
  className
}: {
  title?;
  router;
  fragment;
  queries;
  children;
  onClick?;
  className?;
}) {
  return (
    <a
      className={className || 'tea-text-overflow tea-d-block'}
      title={title}
      href={router.buildUrl(fragment, queries)}
      onClick={event => {
        if (
          event.ctrlKey ||
          event.shiftKey ||
          event.metaKey || // apple
          (event.button && event.button === 1) // middle click, >IE9 + everyone else
        ) {
          return;
        }
        event.preventDefault();
        onClick && onClick(event);
        router.navigate(fragment, queries);
      }}
    >
      {children}
    </a>
  );
}
