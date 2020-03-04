import * as classnames from 'classnames';
import * as React from 'react';

import { SearchBox, SearchBoxProps } from '@tea/component';
import {
    BaseReactProps, findById, findChildren, Identifiable, selectionInsert, selectionRemove, uuid
} from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface ResourceSelectorProps<TResource extends Identifiable> extends BaseReactProps {
  /**
   * 要显示的资源列表
   */
  list: TResource[];

  /*
   * 已选择的资源列表
   **/
  selection?: TResource[];

  /**
   * 如何渲染某一项资源的名称
   */
  itemNameRender: (item: TResource) => string | JSX.Element;

  /**
   * 如何渲染某一项资源的额外信息
   **/
  itemDescriptionRender?: (item: TResource) => string | JSX.Element;

  /**
   * 提供函数判断具体某一项是否被禁用
   */
  itemDisabled?: (item: TResource) => string | boolean;

  /**
   * 用户选择发生改变时回调，此时应该去更新 selection 的数据
   */
  onSelectionChanged: (selection: TResource[]) => void;

  /**
   * 搜索框的配置，具体可参考 `SearchBox` 组件的配置
   **/
  search?: SearchBoxProps;

  /**
   * 选择器的标题，默认为「请选择」
   **/
  selectorTitle?: string;
}

export class ResourceSelectorGeneric<T extends Identifiable> extends React.Component<ResourceSelectorProps<T>, {}> {
  private _secret = uuid();

  public render() {
    const { selectorTitle = t('请选择'), children, search, list, selection } = this.props;

    const heads = findChildren(this.props.children, ResourceSelectorHead);
    const infoRows = findChildren(this.props.children, ResourceSelectorInfoRow);

    return (
      <div className={classnames('tc-15-mod-selector', this.props.className)}>
        {heads.map((head, index) => (
          <div key={index} className="tc-15-mod-selector-area">
            {head}
          </div>
        ))}

        <div className="tc-15-mod-selector-tb">
          {/* left */}
          <div className="tc-15-option-cell options-left">
            <div className="tc-15-option-hd">
              <h4>{selectorTitle}</h4>
            </div>
            <div className="tc-15-option-bd">
              {search && <SearchBox {...search} />}
              <div className="tc-15-option-box">
                {infoRows.map((row, index) => (
                  <div key={index} className="info-row">
                    {row}
                  </div>
                ))}
                <ul className="tc-15-option-list">{list.map((x, index) => this._renderItem(x, index))}</ul>
              </div>
            </div>
          </div>

          {/* seperator */}
          <div className="tc-15-option-cell separator-cell">
            <i className="icon-sep" />
          </div>

          {/* right */}
          <div className="tc-15-option-cell options-right">
            <div className="tc-15-option-hd">
              <h4>
                {t('已选择')}({selection.length})
              </h4>
            </div>
            <div className="tc-15-option-bd">
              <div className="tc-15-option-box">
                {!selection.length && <ResourceSelectorInfoRow>{t('暂未选择')}</ResourceSelectorInfoRow>}
                <ul className="tc-15-option-list">{selection.map((x, index) => this._renderSelectedItem(x, index))}</ul>
              </div>
            </div>
          </div>
        </div>

        {/*<div className="tc-15-mod-selector-tips">支持按住 Shift 键进行多选</div>*/}
      </div>
    );
  }

  private _renderItem(item: T, index) {
    const selected = !!findById(this.props.selection, item.id);
    let disabledReason: any = this.props.itemDisabled && this.props.itemDisabled(item);

    if (disabledReason === true) {
      disabledReason = t('该项不可用');
    }

    return (
      <li key={index} className={classnames({ selected })} title={disabledReason || null}>
        <input
          type="checkbox"
          checked={selected}
          className="tc-15-checkbox"
          id={this._secret + item.id}
          onChange={e => {
            if (!disabledReason) {
              this._checkItem(item, e.target.checked);
            }
          }}
          disabled={!!disabledReason}
        />
        <label className="opt-txt" htmlFor={this._secret + item.id}>
          {this._renderText(item)}
        </label>
      </li>
    );
  }

  private _renderSelectedItem(item: T, index) {
    return (
      <li key={index}>
        <span className="opt-txt">{this._renderText(item)}</span>
        <a
          href="javascript: void(0)"
          className="opt-act"
          role="button"
          title={t('取消')}
          onClick={() => this._checkItem(item, false)}
        >
          <i className="icon-del">{t('取消')}</i>
        </a>
      </li>
    );
  }

  private _renderText(item: T) {
    const { itemNameRender, itemDescriptionRender } = this.props;
    const name = itemNameRender(item);
    const desc = itemDescriptionRender && itemDescriptionRender(item);
    return (
      <span className="opt-txt-inner">
        <span className="item-name">{name}</span>
        {desc && <span className="item-descr">{desc}</span>}
      </span>
    );
  }

  private _checkItem(item: T, checked: boolean) {
    const operation = checked ? selectionInsert : selectionRemove;
    const nextSelection = operation(this.props.list, this.props.selection, item);

    this.props.onSelectionChanged(nextSelection);
  }
}

export function ResourceSelectorHead({ children }: BaseReactProps) {
  return <div className="tc-15-mod-selector-area">{children}</div>;
}

export function ResourceSelectorInfoRow({ children }: BaseReactProps) {
  return <div className="info-row">{children}</div>;
}
