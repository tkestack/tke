import * as React from 'react';

export const Loading = (props) => {

  return (
    <div className="tc-15-autocomplete" style={{left: `${props.offset}px`}}>
      <ul className="tc-15-autocomplete-menu" role="menu">
        <li role="presentation">
          <a className="autocomplete-empty" role="menuitem" href="javascript:;">加载中 ..</a>
        </li>
      </ul>
    </div>
  )
}