/**
 * hpa详情，基本信息tab
 */
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
// import { router } from '@src/modules/cluster/router.project';
import {
  Layout,
  Text,
} from '@tencent/tea-component';
import { MetricsResourceMap } from '../constant';

const Basic = React.memo((props: {
  selectedHpa: any;
}) => {
  // const route = useSelector((state) => state.route);
  // const urlParams = router.resolve(route);
  const { selectedHpa } = props;

  return !isEmpty(selectedHpa) && (
    <ul className="item-descr-list">
      <li>
        <span className="item-descr-tit"><Trans>名称</Trans></span>
        <span className="item-descr-txt">{selectedHpa.metadata.name}</span>
      </li>
      <li>
        <span className="item-descr-tit">Namespace</span>
        <span className="item-descr-txt">{selectedHpa.metadata.namespace}</span>
      </li>
      <li>
        <span className="item-descr-tit"><Trans>关联deployment</Trans></span>
        <span className="item-descr-txt">{selectedHpa.spec.scaleTargetRef.kind}:{selectedHpa.spec.scaleTargetRef.name}</span>
      </li>
      <li>
        <span className="item-descr-tit"><Trans>HPA指标</Trans></span>
        <span className="item-descr-txt">
          {selectedHpa.spec.metrics.map((item, index) => {
            const { name, targetAverageValue, targetAverageUtilization } = item.resource;
            const { meaning, unit } = MetricsResourceMap[name];
            const content = targetAverageValue ? meaning + targetAverageValue : meaning + targetAverageUtilization + unit;
            return <Text key={index} parent="div">{content}</Text>;
          })}
        </span>
      </li>
      <li>
        <span className="item-descr-tit"><Trans>最小副本数</Trans></span>
        <span className="item-descr-txt">{selectedHpa.spec.minReplicas}</span>
      </li>
      <li>
        <span className="item-descr-tit"><Trans>最大副本数</Trans></span>
        <span className="item-descr-txt">{selectedHpa.spec.maxReplicas}</span>
      </li>
    </ul>
  );
});
export default Basic;
